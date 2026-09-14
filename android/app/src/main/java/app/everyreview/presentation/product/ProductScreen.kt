package app.everyreview.presentation.product

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.Button
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.lifecycle.Lifecycle
import androidx.lifecycle.compose.LifecycleEventEffect
import app.everyreview.domain.model.Review

@Composable
fun ProductScreen(
    onWriteReview: (productId: String) -> Unit,
    viewModel: ProductViewModel = hiltViewModel(),
) {
    val uiState by viewModel.uiState.collectAsState()

    LifecycleEventEffect(Lifecycle.Event.ON_RESUME) {
        viewModel.refreshReviews()
    }

    when (val state = uiState) {
        is ProductUiState.Loading -> Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
            CircularProgressIndicator()
        }

        is ProductUiState.Error -> Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
            Text(state.message)
        }

        is ProductUiState.Success -> LazyColumn(
            modifier = Modifier.fillMaxSize(),
            contentPadding = PaddingValues(16.dp),
            verticalArrangement = Arrangement.spacedBy(8.dp),
        ) {
            item {
                Column {
                    Text(
                        text = state.product.name ?: "Unknown product",
                        style = MaterialTheme.typography.headlineSmall,
                    )
                    state.product.brand?.let { Text(it) }
                    Text("Barcode: ${state.product.barcode}")
                    Text(ratingSummaryText(state.product.reviewCount, state.product.averageRating))
                    Button(
                        onClick = { onWriteReview(state.product.id) },
                        modifier = Modifier.fillMaxWidth().padding(top = 8.dp),
                    ) {
                        Text("Write a review")
                    }
                    HorizontalDivider(modifier = Modifier.padding(top = 16.dp))
                }
            }
            if (state.reviews.isEmpty()) {
                item { Text("Be the first to review this product.") }
            } else {
                items(state.reviews, key = { it.id }) { review -> ReviewRow(review) }
            }
        }
    }
}

private fun ratingSummaryText(reviewCount: Int, averageRating: Double?): String {
    if (reviewCount == 0) return "No reviews yet"
    val reviewsWord = if (reviewCount == 1) "review" else "reviews"
    return "★ %.1f (%d %s)".format(averageRating ?: 0.0, reviewCount, reviewsWord)
}

@Composable
private fun ReviewRow(review: Review) {
    Column(modifier = Modifier.fillMaxWidth().padding(vertical = 4.dp)) {
        val ratingText = review.rating?.let { "★".repeat(it) }.orEmpty()
        Text(text = "${review.authorName}  $ratingText", style = MaterialTheme.typography.titleSmall)
        Text(review.body)
    }
}
