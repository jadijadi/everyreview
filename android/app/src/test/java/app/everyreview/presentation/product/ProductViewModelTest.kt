package app.everyreview.presentation.product

import androidx.lifecycle.SavedStateHandle
import app.everyreview.domain.model.Product
import app.everyreview.domain.model.Review
import app.everyreview.domain.repository.ProductRepository
import app.everyreview.domain.repository.ReviewRepository
import app.everyreview.domain.usecase.GetProductUseCase
import app.everyreview.domain.usecase.GetReviewsUseCase
import com.google.common.truth.Truth.assertThat
import io.mockk.coEvery
import io.mockk.mockk
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.StandardTestDispatcher
import kotlinx.coroutines.test.resetMain
import kotlinx.coroutines.test.runTest
import kotlinx.coroutines.test.setMain
import org.junit.After
import org.junit.Before
import org.junit.Test

@OptIn(ExperimentalCoroutinesApi::class)
class ProductViewModelTest {

    private val dispatcher = StandardTestDispatcher()
    private val productRepository = mockk<ProductRepository>()
    private val reviewRepository = mockk<ReviewRepository>()

    @Before
    fun setUp() {
        Dispatchers.setMain(dispatcher)
    }

    @After
    fun tearDown() {
        Dispatchers.resetMain()
    }

    private fun createViewModel(barcode: String = "0000000000000") = ProductViewModel(
        SavedStateHandle(mapOf("barcode" to barcode)),
        GetProductUseCase(productRepository),
        GetReviewsUseCase(reviewRepository),
    )

    @Test
    fun `loads product and reviews on init`() = runTest {
        val product = Product("id-1", "0000000000000", "Widget", null, null, false, 1, 5.0)
        val review = Review("r1", "id-1", "Alice", "Nice", 5, "2026-01-01T00:00:00Z")
        coEvery { productRepository.getByBarcode("0000000000000") } returns Result.success(product)
        coEvery { reviewRepository.listByProduct("id-1") } returns Result.success(listOf(review))

        val viewModel = createViewModel()
        dispatcher.scheduler.advanceUntilIdle()

        val state = viewModel.uiState.value
        assertThat(state).isInstanceOf(ProductUiState.Success::class.java)
        state as ProductUiState.Success
        assertThat(state.product).isEqualTo(product)
        assertThat(state.reviews).containsExactly(review)
    }

    @Test
    fun `surfaces error when product lookup fails`() = runTest {
        coEvery { productRepository.getByBarcode(any()) } returns Result.failure(RuntimeException("offline"))

        val viewModel = createViewModel()
        dispatcher.scheduler.advanceUntilIdle()

        assertThat(viewModel.uiState.value).isInstanceOf(ProductUiState.Error::class.java)
    }
}
