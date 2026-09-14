package app.everyreview.presentation.product

import app.everyreview.domain.model.Product
import app.everyreview.domain.model.Review

sealed interface ProductUiState {
    data object Loading : ProductUiState

    data class Error(val message: String) : ProductUiState

    data class Success(
        val product: Product,
        val reviews: List<Review>,
    ) : ProductUiState
}
