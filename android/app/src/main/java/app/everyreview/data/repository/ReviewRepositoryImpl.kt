package app.everyreview.data.repository

import app.everyreview.data.local.ReviewDao
import app.everyreview.data.remote.EveryReviewApi
import app.everyreview.data.remote.dto.CreateReviewRequest
import app.everyreview.domain.model.Review
import app.everyreview.domain.repository.ReviewRepository
import kotlinx.coroutines.CancellationException
import javax.inject.Inject

class ReviewRepositoryImpl @Inject constructor(
    private val api: EveryReviewApi,
    private val reviewDao: ReviewDao,
) : ReviewRepository {

    override suspend fun listByProduct(productId: String): Result<List<Review>> {
        return try {
            val reviews = api.getReviews(productId).items.map { it.toDomain() }
            reviewDao.deleteByProduct(productId)
            reviewDao.upsertAll(reviews.map { it.toEntity() })
            Result.success(reviews)
        } catch (e: CancellationException) {
            throw e
        } catch (e: Exception) {
            val cached = reviewDao.listByProduct(productId).map { it.toDomain() }
            if (cached.isNotEmpty()) Result.success(cached) else Result.failure(e)
        }
    }

    override suspend fun submit(
        productId: String,
        authorName: String,
        body: String,
        rating: Int?,
    ): Result<Review> {
        return try {
            val created = api.createReview(
                productId,
                CreateReviewRequest(authorName = authorName, body = body, rating = rating),
            ).toDomain()
            Result.success(created)
        } catch (e: CancellationException) {
            throw e
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
}
