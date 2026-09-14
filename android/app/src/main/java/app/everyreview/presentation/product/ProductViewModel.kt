package app.everyreview.presentation.product

import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import app.everyreview.domain.usecase.GetProductUseCase
import app.everyreview.domain.usecase.GetReviewsUseCase
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class ProductViewModel @Inject constructor(
    savedStateHandle: SavedStateHandle,
    private val getProduct: GetProductUseCase,
    private val getReviews: GetReviewsUseCase,
) : ViewModel() {

    private val barcode: String = checkNotNull(savedStateHandle["barcode"])

    private val _uiState = MutableStateFlow<ProductUiState>(ProductUiState.Loading)
    val uiState: StateFlow<ProductUiState> = _uiState.asStateFlow()

    init {
        refresh()
    }

    fun refresh() {
        viewModelScope.launch {
            _uiState.value = ProductUiState.Loading
            getProduct(barcode)
                .onSuccess { product ->
                    getReviews(product.id)
                        .onSuccess { reviews -> _uiState.value = ProductUiState.Success(product, reviews) }
                        .onFailure { _uiState.value = ProductUiState.Success(product, emptyList()) }
                }
                .onFailure { error ->
                    _uiState.value = ProductUiState.Error(error.message ?: "Something went wrong")
                }
        }
    }

    /**
     * Re-fetches product and reviews while keeping the current content on screen (used after
     * popping back from Write Review / Add Details). Failures keep whatever was shown before.
     */
    fun refreshQuietly() {
        val current = _uiState.value
        if (current !is ProductUiState.Success) return
        viewModelScope.launch {
            val product = getProduct(barcode).getOrDefault(current.product)
            val reviews = getReviews(product.id).getOrDefault(current.reviews)
            _uiState.value = ProductUiState.Success(product, reviews)
        }
    }
}
