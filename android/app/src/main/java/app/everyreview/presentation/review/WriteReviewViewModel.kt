package app.everyreview.presentation.review

import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import app.everyreview.data.prefs.NicknamePrefs
import app.everyreview.domain.usecase.SubmitReviewUseCase
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import javax.inject.Inject

data class WriteReviewUiState(
    val nickname: String = "",
    val body: String = "",
    val rating: Int? = null,
    val isSubmitting: Boolean = false,
    val error: String? = null,
    val submitted: Boolean = false,
)

@HiltViewModel
class WriteReviewViewModel @Inject constructor(
    savedStateHandle: SavedStateHandle,
    private val submitReview: SubmitReviewUseCase,
    private val nicknamePrefs: NicknamePrefs,
) : ViewModel() {

    private val productId: String = checkNotNull(savedStateHandle["productId"])

    private val _uiState = MutableStateFlow(WriteReviewUiState())
    val uiState: StateFlow<WriteReviewUiState> = _uiState.asStateFlow()

    init {
        viewModelScope.launch {
            _uiState.update { it.copy(nickname = nicknamePrefs.getNickname()) }
        }
    }

    fun onNicknameChanged(value: String) = _uiState.update { it.copy(nickname = value) }

    fun onBodyChanged(value: String) = _uiState.update { it.copy(body = value, error = null) }

    fun onRatingChanged(value: Int) = _uiState.update { it.copy(rating = value) }

    fun submit() {
        val state = _uiState.value
        viewModelScope.launch {
            _uiState.update { it.copy(isSubmitting = true, error = null) }
            nicknamePrefs.setNickname(state.nickname)
            submitReview(productId, state.nickname, state.body, state.rating)
                .onSuccess {
                    _uiState.update { it.copy(isSubmitting = false, submitted = true) }
                }
                .onFailure { e ->
                    _uiState.update { it.copy(isSubmitting = false, error = e.message ?: "Could not submit review") }
                }
        }
    }
}
