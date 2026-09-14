package app.everyreview.presentation.addproduct

import androidx.lifecycle.SavedStateHandle
import app.everyreview.data.media.PhotoPreparer
import app.everyreview.domain.model.Product
import app.everyreview.domain.model.ProductDetails
import app.everyreview.domain.repository.ProductRepository
import app.everyreview.domain.usecase.SubmitProductDetailsUseCase
import com.google.common.truth.Truth.assertThat
import io.mockk.coEvery
import io.mockk.coVerify
import io.mockk.mockk
import io.mockk.slot
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
class AddProductViewModelTest {

    private val dispatcher = StandardTestDispatcher()
    private val productRepository = mockk<ProductRepository>()
    private val photoPreparer = mockk<PhotoPreparer>()

    @Before
    fun setUp() {
        Dispatchers.setMain(dispatcher)
    }

    @After
    fun tearDown() {
        Dispatchers.resetMain()
    }

    private fun createViewModel() = AddProductViewModel(
        SavedStateHandle(mapOf("productId" to "product-1", "barcode" to "4006381333931")),
        SubmitProductDetailsUseCase(productRepository),
        photoPreparer,
    )

    @Test
    fun `starts with the barcode and cannot submit until a name is typed`() {
        val viewModel = createViewModel()

        assertThat(viewModel.uiState.value.barcode).isEqualTo("4006381333931")
        assertThat(viewModel.uiState.value.canSubmit).isFalse()

        viewModel.onNameChanged("Highlighter")

        assertThat(viewModel.uiState.value.canSubmit).isTrue()
    }

    @Test
    fun `submit sends the form and flags completion`() = runTest {
        val saved = Product("product-1", "4006381333931", "Highlighter", "Stabilo", null, false, 0, null)
        val sent = slot<ProductDetails>()
        coEvery { productRepository.submitDetails("product-1", capture(sent)) } returns Result.success(saved)
        val viewModel = createViewModel()
        viewModel.onNameChanged("Highlighter")
        viewModel.onBrandChanged("Stabilo")
        viewModel.onPriceChanged("2.49")
        viewModel.onCurrencyChanged("EUR")

        viewModel.submit()
        dispatcher.scheduler.advanceUntilIdle()

        assertThat(viewModel.uiState.value.submitted).isTrue()
        assertThat(viewModel.uiState.value.isSubmitting).isFalse()
        assertThat(sent.captured.name).isEqualTo("Highlighter")
        assertThat(sent.captured.brand).isEqualTo("Stabilo")
        assertThat(sent.captured.price).isEqualTo("2.49")
        assertThat(sent.captured.currency).isEqualTo("EUR")
        assertThat(sent.captured.photo).isNull()
    }

    @Test
    fun `submit surfaces the failure and stays on the form`() = runTest {
        coEvery { productRepository.submitDetails(any(), any()) } returns
            Result.failure(RuntimeException("Someone already added details for this product."))
        val viewModel = createViewModel()
        viewModel.onNameChanged("Highlighter")

        viewModel.submit()
        dispatcher.scheduler.advanceUntilIdle()

        assertThat(viewModel.uiState.value.submitted).isFalse()
        assertThat(viewModel.uiState.value.error).contains("already added")
    }

    @Test
    fun `submit is a no-op while the name is blank`() = runTest {
        val viewModel = createViewModel()

        viewModel.submit()
        dispatcher.scheduler.advanceUntilIdle()

        coVerify(exactly = 0) { productRepository.submitDetails(any(), any()) }
    }
}
