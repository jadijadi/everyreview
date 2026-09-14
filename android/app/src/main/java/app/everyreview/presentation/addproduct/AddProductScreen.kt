package app.everyreview.presentation.addproduct

import android.Manifest
import android.content.pm.PackageManager
import android.net.Uri
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.PickVisualMediaRequest
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.imePadding
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.input.KeyboardCapitalization
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.dp
import androidx.core.content.ContextCompat
import androidx.hilt.navigation.compose.hiltViewModel
import coil.compose.AsyncImage

@Composable
fun AddProductScreen(
    onSubmitted: () -> Unit,
    viewModel: AddProductViewModel = hiltViewModel(),
) {
    val uiState by viewModel.uiState.collectAsState()

    LaunchedEffect(uiState.submitted) {
        if (uiState.submitted) onSubmitted()
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .verticalScroll(rememberScrollState())
            .imePadding()
            .padding(16.dp),
    ) {
        Text("Add product details", style = MaterialTheme.typography.headlineSmall)
        Text(
            text = "Barcode ${uiState.barcode} isn't in EveryReview yet. Tell others what it is.",
            style = MaterialTheme.typography.bodyMedium,
            modifier = Modifier.padding(top = 4.dp),
        )

        PhotoPicker(
            photoUri = uiState.photoUri,
            newCaptureUri = viewModel::newCaptureUri,
            onPhotoPicked = viewModel::onPhotoPicked,
            onPhotoRemoved = viewModel::onPhotoRemoved,
            modifier = Modifier.padding(top = 16.dp),
        )

        OutlinedTextField(
            value = uiState.name,
            onValueChange = viewModel::onNameChanged,
            label = { Text("Product name *") },
            singleLine = true,
            keyboardOptions = KeyboardOptions(capitalization = KeyboardCapitalization.Sentences),
            modifier = Modifier.fillMaxWidth().padding(top = 16.dp),
        )
        OutlinedTextField(
            value = uiState.brand,
            onValueChange = viewModel::onBrandChanged,
            label = { Text("Brand") },
            singleLine = true,
            keyboardOptions = KeyboardOptions(capitalization = KeyboardCapitalization.Words),
            modifier = Modifier.fillMaxWidth().padding(top = 8.dp),
        )
        OutlinedTextField(
            value = uiState.manufacturer,
            onValueChange = viewModel::onManufacturerChanged,
            label = { Text("Manufacturer") },
            singleLine = true,
            keyboardOptions = KeyboardOptions(capitalization = KeyboardCapitalization.Words),
            modifier = Modifier.fillMaxWidth().padding(top = 8.dp),
        )
        Row(modifier = Modifier.fillMaxWidth().padding(top = 8.dp), horizontalArrangement = Arrangement.spacedBy(8.dp)) {
            OutlinedTextField(
                value = uiState.price,
                onValueChange = viewModel::onPriceChanged,
                label = { Text("Price") },
                singleLine = true,
                keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Decimal),
                modifier = Modifier.weight(1f),
            )
            OutlinedTextField(
                value = uiState.currency,
                onValueChange = viewModel::onCurrencyChanged,
                label = { Text("Currency") },
                singleLine = true,
                keyboardOptions = KeyboardOptions(capitalization = KeyboardCapitalization.Characters),
                modifier = Modifier.width(120.dp),
            )
        }
        OutlinedTextField(
            value = uiState.description,
            onValueChange = viewModel::onDescriptionChanged,
            label = { Text("Description") },
            keyboardOptions = KeyboardOptions(capitalization = KeyboardCapitalization.Sentences),
            modifier = Modifier.fillMaxWidth().padding(top = 8.dp).height(120.dp),
        )

        uiState.error?.let { message ->
            Text(
                text = message,
                color = MaterialTheme.colorScheme.error,
                modifier = Modifier.padding(top = 8.dp),
            )
        }

        Button(
            onClick = viewModel::submit,
            enabled = uiState.canSubmit,
            modifier = Modifier.fillMaxWidth().padding(top = 16.dp),
        ) {
            Text(if (uiState.isSubmitting) "Saving…" else "Save product")
        }
    }
}

/**
 * Camera or gallery. The app declares CAMERA (for the scanner), so Android requires us to
 * hold it before launching the camera intent too — hence the permission check here.
 */
@Composable
private fun PhotoPicker(
    photoUri: Uri?,
    newCaptureUri: () -> Uri,
    onPhotoPicked: (Uri?) -> Unit,
    onPhotoRemoved: () -> Unit,
    modifier: Modifier = Modifier,
) {
    val context = LocalContext.current
    var pendingCaptureUri by rememberSaveable { mutableStateOf<Uri?>(null) }

    val takePicture = rememberLauncherForActivityResult(ActivityResultContracts.TakePicture()) { saved ->
        if (saved) onPhotoPicked(pendingCaptureUri)
    }
    val launchCamera = {
        val uri = newCaptureUri()
        pendingCaptureUri = uri
        takePicture.launch(uri)
    }
    val requestCamera = rememberLauncherForActivityResult(ActivityResultContracts.RequestPermission()) { granted ->
        if (granted) launchCamera()
    }
    val pickFromGallery = rememberLauncherForActivityResult(ActivityResultContracts.PickVisualMedia()) { uri ->
        if (uri != null) onPhotoPicked(uri)
    }

    Column(modifier = modifier) {
        if (photoUri != null) {
            AsyncImage(
                model = photoUri,
                contentDescription = "Product photo",
                contentScale = ContentScale.Fit,
                modifier = Modifier.fillMaxWidth().height(200.dp),
            )
            TextButton(onClick = onPhotoRemoved, modifier = Modifier.align(Alignment.End)) { Text("Remove photo") }
        } else {
            Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                OutlinedButton(
                    onClick = {
                        val granted = ContextCompat.checkSelfPermission(context, Manifest.permission.CAMERA) ==
                            PackageManager.PERMISSION_GRANTED
                        if (granted) launchCamera() else requestCamera.launch(Manifest.permission.CAMERA)
                    },
                    modifier = Modifier.weight(1f),
                ) { Text("Take photo") }
                OutlinedButton(
                    onClick = {
                        pickFromGallery.launch(PickVisualMediaRequest(ActivityResultContracts.PickVisualMedia.ImageOnly))
                    },
                    modifier = Modifier.weight(1f),
                ) { Text("Choose photo") }
            }
        }
    }
}
