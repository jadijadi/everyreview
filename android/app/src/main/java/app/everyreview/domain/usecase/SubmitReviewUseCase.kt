package app.everyreview.domain.usecase

import app.everyreview.domain.model.Review
import app.everyreview.domain.repository.ReviewRepository
import javax.inject.Inject

class SubmitReviewUseCase @Inject constructor(
    private val repository: ReviewRepository,
) {
    suspend operator fun invoke(
        productId: String,
        authorName: String,
        body: String,
        rating: Int?,
    ): Result<Review> {
        val trimmedBody = body.trim()
        if (trimmedBody.isEmpty()) {
            return Result.failure(IllegalArgumentException("Review body must not be empty"))
        }
        if (rating != null && rating !in 1..5) {
            return Result.failure(IllegalArgumentException("Rating must be between 1 and 5"))
        }
        return repository.submit(productId, authorName.trim(), trimmedBody, rating)
    }
}
