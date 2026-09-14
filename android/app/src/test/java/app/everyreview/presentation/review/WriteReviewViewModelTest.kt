package app.everyreview.presentation.review

import androidx.lifecycle.SavedStateHandle
import app.everyreview.data.prefs.NicknamePrefs
import app.everyreview.domain.model.Review
import app.everyreview.domain.repository.ReviewRepository
import app.everyreview.domain.usecase.SubmitReviewUseCase
import com.google.common.truth.Truth.assertThat
import io.mockk.coEvery
import io.mockk.coJustRun
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
class WriteReviewViewModelTest {

    private val dispatcher = StandardTestDispatcher()
    private val reviewRepository = mockk<ReviewRepository>()
    private val nicknamePrefs = mockk<NicknamePrefs>()

    @Before
    fun setUp() {
        Dispatchers.setMain(dispatcher)
        coEvery { nicknamePrefs.getNickname() } returns "Alice"
        coJustRun { nicknamePrefs.setNickname(any()) }
    }

    @After
    fun tearDown() {
        Dispatchers.resetMain()
    }

    private fun createViewModel() = WriteReviewViewModel(
        SavedStateHandle(mapOf("productId" to "product-1")),
        SubmitReviewUseCase(reviewRepository),
        nicknamePrefs,
    )

    @Test
    fun `preloads saved nickname`() = runTest {
        val viewModel = createViewModel()
        dispatcher.scheduler.advanceUntilIdle()

        assertThat(viewModel.uiState.value.nickname).isEqualTo("Alice")
    }

    @Test
    fun `submit failure surfaces error and does not mark submitted`() = runTest {
        val viewModel = createViewModel()
        dispatcher.scheduler.advanceUntilIdle()
        viewModel.onBodyChanged("")

        viewModel.submit()
        dispatcher.scheduler.advanceUntilIdle()

        val state = viewModel.uiState.value
        assertThat(state.submitted).isFalse()
        assertThat(state.error).isNotNull()
    }

    @Test
    fun `successful submit marks state as submitted`() = runTest {
        val created = Review("r1", "product-1", "Alice", "Great!", 5, "2026-01-01T00:00:00Z")
        coEvery { reviewRepository.submit("product-1", "Alice", "Great!", 5) } returns Result.success(created)

        val viewModel = createViewModel()
        dispatcher.scheduler.advanceUntilIdle()
        viewModel.onBodyChanged("Great!")
        viewModel.onRatingChanged(5)

        viewModel.submit()
        dispatcher.scheduler.advanceUntilIdle()

        assertThat(viewModel.uiState.value.submitted).isTrue()
    }
}
