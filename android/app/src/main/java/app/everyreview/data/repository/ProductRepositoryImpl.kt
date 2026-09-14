package app.everyreview.data.repository

import app.everyreview.data.local.ProductDao
import app.everyreview.data.remote.EveryReviewApi
import app.everyreview.domain.model.Product
import app.everyreview.domain.repository.ProductRepository
import kotlinx.coroutines.CancellationException
import javax.inject.Inject

class ProductRepositoryImpl @Inject constructor(
    private val api: EveryReviewApi,
    private val productDao: ProductDao,
) : ProductRepository {

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
}
