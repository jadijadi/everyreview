package app.everyreview.data.remote.dto

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class ProductResponse(
    val id: String,
    val barcode: String,
    val name: String?,
    val brand: String?,
    @Json(name = "image_url") val imageUrl: String?,
    @Json(name = "is_placeholder") val isPlaceholder: Boolean,
    @Json(name = "rating_summary") val ratingSummary: RatingSummaryResponse,
)

@JsonClass(generateAdapter = true)
data class RatingSummaryResponse(
    @Json(name = "review_count") val reviewCount: Int,
    @Json(name = "average_rating") val averageRating: Double?,
)
