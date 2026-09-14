package app.everyreview.domain.usecase

import app.everyreview.domain.model.Review
import app.everyreview.domain.repository.ReviewRepository
import com.google.common.truth.Truth.assertThat
import io.mockk.coEvery
import io.mockk.coVerify
import io.mockk.mockk
import kotlinx.coroutines.test.runTest
import org.junit.Test

class SubmitReviewUseCaseTest {

    private val repository = mockk<ReviewRepository>()
    private val useCase = SubmitReviewUseCase(repository)

    @Test
    fun `rejects blank body without calling repository`() = runTest {
        val result = useCase("product-1", "Alice", "   ", null)

        assertThat(result.isFailure).isTrue()
        coVerify(exactly = 0) { repository.submit(any(), any(), any(), any()) }
    }

    @Test
    fun `rejects rating outside 1 to 5`() = runTest {
        val result = useCase("product-1", "Alice", "Great!", 6)

        assertThat(result.isFailure).isTrue()
    }

    @Test
    fun `trims body and author before delegating to repository`() = runTest {
        val expected = Review("r1", "product-1", "Alice", "Great!", 5, "2026-01-01T00:00:00Z")
        coEvery { repository.submit("product-1", "Alice", "Great!", 5) } returns Result.success(expected)

        val result = useCase("product-1", "  Alice  ", "  Great!  ", 5)

        assertThat(result.getOrNull()).isEqualTo(expected)
        coVerify { repository.submit("product-1", "Alice", "Great!", 5) }
    }
}
