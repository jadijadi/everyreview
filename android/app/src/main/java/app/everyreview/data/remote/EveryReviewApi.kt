package app.everyreview.data.remote

import app.everyreview.data.remote.dto.CreateReviewRequest
import app.everyreview.data.remote.dto.ProductResponse
import app.everyreview.data.remote.dto.ReviewPageResponse
import app.everyreview.data.remote.dto.ReviewResponse
import okhttp3.MultipartBody
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.Multipart
import retrofit2.http.POST
import retrofit2.http.Part
import retrofit2.http.Path

/**
 * Hand-written against the trimmed MVP backend (backend/README.md "MVP scope"),
 * not generated from backend/api/openapi.yaml — see that file's scope notes.
 */
interface EveryReviewApi {
    @GET("products/{barcode}")
    suspend fun getProduct(@Path("barcode") barcode: String): ProductResponse

    /** Multipart so the optional photo travels with the text fields; see the spec for field names. */
    @Multipart
    @POST("products/{productId}/details")
    suspend fun submitProductDetails(
        @Path("productId") productId: String,
        @Part parts: List<MultipartBody.Part>,
    ): ProductResponse

    @GET("products/{productId}/reviews")
    suspend fun getReviews(@Path("productId") productId: String): ReviewPageResponse

    @POST("products/{productId}/reviews")
    suspend fun createReview(
        @Path("productId") productId: String,
        @Body request: CreateReviewRequest,
    ): ReviewResponse
}
