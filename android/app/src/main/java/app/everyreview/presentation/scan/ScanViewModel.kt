package app.everyreview.presentation.scan

import androidx.lifecycle.ViewModel
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import javax.inject.Inject

/**
 * Guards against acting on more than one barcode per visit to this screen (ML Kit reports many
 * frames per second). The ViewModel outlives a trip to the Product screen — the scan destination
 * stays on the back stack — so the screen must call [reset] each time it comes back into view.
 */
@HiltViewModel
class ScanViewModel @Inject constructor() : ViewModel() {
    private val _hasScanned = MutableStateFlow(false)
    val hasScanned: StateFlow<Boolean> = _hasScanned.asStateFlow()

    fun onBarcodeDetected(): Boolean {
        if (_hasScanned.value) return false
        _hasScanned.value = true
        return true
    }

    fun reset() {
        _hasScanned.value = false
    }
}
