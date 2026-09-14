package app.everyreview.domain.model

data class Review(
    val id: String,
    val productId: String,
    val authorName: String,
    val body: String,
    val rating: Int?,
    val createdAt: String,
)
