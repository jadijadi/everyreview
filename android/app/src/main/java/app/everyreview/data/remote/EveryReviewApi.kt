package app.everyreview.data.remote

import app.everyreview.data.remote.dto.CreateReviewRequest
import app.everyreview.data.remote.dto.ProductResponse
import app.everyreview.data.remote.dto.ReviewPageResponse
import app.everyreview.data.remote.dto.ReviewResponse
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.Path

/**
 * Hand-written against the trimmed MVP backend (backend/README.md "MVP scope"),
 * not generated from backend/api/openapi.yaml — see that file's scope notes.
 */
interface EveryReviewApi {
    @GET("products/{barcode}")
    suspend fun getProduct(@Path("barcode") barcode: String): ProductResponse

    @GET("products/{productId}/reviews")
    suspend fun getReviews(@Path("productId") productId: String): ReviewPageResponse

    @POST("products/{productId}/reviews")
    suspend fun createReview(
        @Path("productId") productId: String,
        @Body request: CreateReviewRequest,
    ): ReviewResponse
}
