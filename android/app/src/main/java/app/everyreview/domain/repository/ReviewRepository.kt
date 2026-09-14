package app.everyreview.domain.repository

import app.everyreview.domain.model.Review

interface ReviewRepository {
    suspend fun listByProduct(productId: String): Result<List<Review>>

    suspend fun submit(
        productId: String,
        authorName: String,
        body: String,
        rating: Int?,
    ): Result<Review>
}
