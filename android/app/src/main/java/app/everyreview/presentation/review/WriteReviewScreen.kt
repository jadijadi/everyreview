package app.everyreview.presentation.review

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel

@Composable
fun WriteReviewScreen(
    onSubmitted: () -> Unit,
    viewModel: WriteReviewViewModel = hiltViewModel(),
) {
    val uiState by viewModel.uiState.collectAsState()

    LaunchedEffect(uiState.submitted) {
        if (uiState.submitted) onSubmitted()
    }

    Column(modifier = Modifier.fillMaxSize().padding(16.dp)) {
        Text("Write a review", style = MaterialTheme.typography.headlineSmall)

        OutlinedTextField(
            value = uiState.nickname,
            onValueChange = viewModel::onNicknameChanged,
            label = { Text("Nickname (optional)") },
            singleLine = true,
            modifier = Modifier.fillMaxWidth().padding(top = 16.dp),
        )

        Row(modifier = Modifier.padding(top = 16.dp)) {
            for (star in 1..5) {
                Text(
                    text = if ((uiState.rating ?: 0) >= star) "★" else "☆",
                    style = MaterialTheme.typography.headlineMedium,
                    modifier = Modifier
                        .clickable { viewModel.onRatingChanged(star) }
                        .padding(end = 4.dp),
                )
            }
        }

        OutlinedTextField(
            value = uiState.body,
            onValueChange = viewModel::onBodyChanged,
            label = { Text("Your review") },
            modifier = Modifier.fillMaxWidth().padding(top = 16.dp).height(140.dp),
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
            enabled = !uiState.isSubmitting,
            modifier = Modifier.fillMaxWidth().padding(top = 16.dp),
        ) {
            Text(if (uiState.isSubmitting) "Submitting…" else "Submit review")
        }
    }
}
