package app.everyreview.domain.usecase

import app.everyreview.domain.model.Review
import app.everyreview.domain.repository.ReviewRepository
import javax.inject.Inject

class GetReviewsUseCase @Inject constructor(
    private val repository: ReviewRepository,
) {
    suspend operator fun invoke(productId: String): Result<List<Review>> = repository.listByProduct(productId)
}
