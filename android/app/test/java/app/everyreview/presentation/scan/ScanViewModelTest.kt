package app.everyreview.presentation.scan

import com.google.common.truth.Truth.assertThat
import org.junit.Test

class ScanViewModelTest {

    @Test
    fun `only the first detection per visit is acted on`() {
        val viewModel = ScanViewModel()

        assertThat(viewModel.onBarcodeDetected()).isTrue()
        assertThat(viewModel.onBarcodeDetected()).isFalse()
        assertThat(viewModel.onBarcodeDetected()).isFalse()
    }

    @Test
    fun `reset re-arms the guard when the screen is shown again`() {
        val viewModel = ScanViewModel()
        viewModel.onBarcodeDetected()

        viewModel.reset()

        assertThat(viewModel.hasScanned.value).isFalse()
        assertThat(viewModel.onBarcodeDetected()).isTrue()
    }
}
