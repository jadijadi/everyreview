package app.everyreview.data.remote.dto

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class ReviewResponse(
    val id: String,
    @Json(name = "product_id") val productId: String,
    @Json(name = "author_name") val authorName: String,
    val body: String,
    val rating: Int?,
    @Json(name = "created_at") val createdAt: String,
)

@JsonClass(generateAdapter = true)
data class ReviewPageResponse(
    val items: List<ReviewResponse>,
)

@JsonClass(generateAdapter = true)
data class CreateReviewRequest(
    @Json(name = "author_name") val authorName: String,
    val body: String,
    val rating: Int?,
)
