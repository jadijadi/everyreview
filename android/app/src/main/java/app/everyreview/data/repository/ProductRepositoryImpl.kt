package app.everyreview.data.repository

import app.everyreview.data.local.ProductDao
import app.everyreview.data.remote.EveryReviewApi
import app.everyreview.data.remote.dto.ApiError
import app.everyreview.domain.model.Product
import app.everyreview.domain.model.ProductDetails
import app.everyreview.domain.repository.ProductRepository
import com.squareup.moshi.JsonAdapter
import com.squareup.moshi.Moshi
import kotlinx.coroutines.CancellationException
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.MultipartBody
import okhttp3.RequestBody.Companion.asRequestBody
import retrofit2.HttpException
import java.io.IOException
import javax.inject.Inject

class ProductRepositoryImpl @Inject constructor(
    private val api: EveryReviewApi,
    private val productDao: ProductDao,
    moshi: Moshi,
) : ProductRepository {

    private val errorAdapter: JsonAdapter<ApiError> = moshi.adapter(ApiError::class.java)

    override suspend fun getByBarcode(barcode: String): Result<Product> {
        return try {
            val product = api.getProduct(barcode).toDomain()
            productDao.upsert(product.toEntity())
            Result.success(product)
        } catch (e: CancellationException) {
            throw e
        } catch (e: Exception) {
            productDao.getByBarcode(barcode)?.let { return Result.success(it.toDomain()) }
            Result.failure(e)
        }
    }

    override suspend fun submitDetails(productId: String, details: ProductDetails): Result<Product> {
        val parts = buildList {
            add(MultipartBody.Part.createFormData("name", details.name))
            details.brand?.let { add(MultipartBody.Part.createFormData("brand", it)) }
            details.manufacturer?.let { add(MultipartBody.Part.createFormData("manufacturer", it)) }
            details.description?.let { add(MultipartBody.Part.createFormData("description", it)) }
            details.price?.let { add(MultipartBody.Part.createFormData("price", it)) }
            details.currency?.let { add(MultipartBody.Part.createFormData("currency", it)) }
            details.photo?.let { file ->
                add(MultipartBody.Part.createFormData("image", file.name, file.asRequestBody("image/jpeg".toMediaType())))
            }
        }
        return try {
            val product = api.submitProductDetails(productId, parts).toDomain()
            productDao.upsert(product.toEntity())
            Result.success(product)
        } catch (e: CancellationException) {
            throw e
        } catch (e: HttpException) {
            Result.failure(IOException(userMessage(e), e))
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /** The backend answers validation failures with `{"error": "..."}`; surface that text instead of "HTTP 400". */
    private fun userMessage(e: HttpException): String {
        val body = e.response()?.errorBody()?.string()
        val parsed = body?.let { runCatching { errorAdapter.fromJson(it) }.getOrNull() }?.error
        return when {
            e.code() == 409 -> "Someone already added details for this product."
            parsed != null -> parsed.removePrefix("product: ").removePrefix("media: ")
            else -> "Could not save product details (HTTP ${e.code()})"
        }
    }
}
