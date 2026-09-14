package app.everyreview.presentation.addproduct

import android.net.Uri
import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import app.everyreview.data.media.PhotoPreparer
import app.everyreview.domain.model.ProductDetails
import app.everyreview.domain.usecase.SubmitProductDetailsUseCase
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import java.util.Currency
import java.util.Locale
import javax.inject.Inject

data class AddProductUiState(
    val barcode: String = "",
    val name: String = "",
    val brand: String = "",
    val manufacturer: String = "",
    val description: String = "",
    val price: String = "",
    val currency: String = "",
    val photoUri: Uri? = null,
    val isSubmitting: Boolean = false,
    val error: String? = null,
    val submitted: Boolean = false,
) {
    val canSubmit: Boolean get() = name.isNotBlank() && !isSubmitting
}

@HiltViewModel
class AddProductViewModel @Inject constructor(
    savedStateHandle: SavedStateHandle,
    private val submitDetails: SubmitProductDetailsUseCase,
    private val photoPreparer: PhotoPreparer,
) : ViewModel() {

    private val productId: String = checkNotNull(savedStateHandle["productId"])

    private val _uiState = MutableStateFlow(
        AddProductUiState(
            barcode = checkNotNull(savedStateHandle["barcode"]),
            currency = defaultCurrencyCode(),
        ),
    )
    val uiState: StateFlow<AddProductUiState> = _uiState.asStateFlow()

    fun onNameChanged(value: String) = _uiState.update { it.copy(name = value, error = null) }

    fun onBrandChanged(value: String) = _uiState.update { it.copy(brand = value) }

    fun onManufacturerChanged(value: String) = _uiState.update { it.copy(manufacturer = value) }

    fun onDescriptionChanged(value: String) = _uiState.update { it.copy(description = value) }

    fun onPriceChanged(value: String) = _uiState.update { it.copy(price = value, error = null) }

    fun onCurrencyChanged(value: String) = _uiState.update { it.copy(currency = value.take(3), error = null) }

    fun onPhotoPicked(uri: Uri?) = _uiState.update { it.copy(photoUri = uri, error = null) }

    fun onPhotoRemoved() = _uiState.update { it.copy(photoUri = null) }

    /** Hands the UI a URI the camera app can write into; the pick is confirmed via [onPhotoPicked]. */
    fun newCaptureUri(): Uri = photoPreparer.newCaptureUri()

    fun submit() {
        val state = _uiState.value
        if (!state.canSubmit) return
        viewModelScope.launch {
            _uiState.update { it.copy(isSubmitting = true, error = null) }
            val photo = state.photoUri?.let { uri ->
                runCatching { photoPreparer.prepareForUpload(uri) }.getOrElse { e ->
                    _uiState.update { it.copy(isSubmitting = false, error = e.message ?: "Could not read the photo") }
                    return@launch
                }
            }
            try {
                submitDetails(
                    productId,
                    ProductDetails(
                        name = state.name,
                        brand = state.brand,
                        manufacturer = state.manufacturer,
                        description = state.description,
                        price = state.price,
                        currency = state.currency,
                        photo = photo,
                    ),
                )
                    .onSuccess { _uiState.update { it.copy(isSubmitting = false, submitted = true) } }
                    .onFailure { e ->
                        _uiState.update { it.copy(isSubmitting = false, error = e.message ?: "Could not save product details") }
                    }
            } finally {
                photo?.delete()
            }
        }
    }

    private fun defaultCurrencyCode(): String =
        runCatching { Currency.getInstance(Locale.getDefault()).currencyCode }.getOrDefault("")
}
