package app.everyreview.presentation.product

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.lifecycle.Lifecycle
import androidx.lifecycle.compose.LifecycleEventEffect
import app.everyreview.domain.model.Product
import app.everyreview.domain.model.Review
import coil.compose.AsyncImage
import java.text.NumberFormat
import java.util.Currency

@Composable
fun ProductScreen(
    onWriteReview: (productId: String) -> Unit,
    onAddDetails: (productId: String, barcode: String) -> Unit,
    viewModel: ProductViewModel = hiltViewModel(),
) {
    val uiState by viewModel.uiState.collectAsState()

    // Coming back from Write Review or Add Details: pick up the new review / details without a loading flash.
    LifecycleEventEffect(Lifecycle.Event.ON_RESUME) {
        viewModel.refreshQuietly()
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
                ProductHeader(
                    product = state.product,
                    onWriteReview = { onWriteReview(state.product.id) },
                    onAddDetails = { onAddDetails(state.product.id, state.product.barcode) },
                )
            }
            if (state.reviews.isEmpty()) {
                item { Text("Be the first to review this product.") }
            } else {
                items(state.reviews, key = { it.id }) { review -> ReviewRow(review) }
            }
        }
    }
}

@Composable
private fun ProductHeader(product: Product, onWriteReview: () -> Unit, onAddDetails: () -> Unit) {
    Column {
        product.imageUrl?.let { url ->
            AsyncImage(
                model = url,
                contentDescription = product.name,
                contentScale = ContentScale.Fit,
                modifier = Modifier.fillMaxWidth().height(220.dp).padding(bottom = 12.dp),
            )
        }
        Text(
            text = product.name ?: "Unknown product",
            style = MaterialTheme.typography.headlineSmall,
        )
        listOfNotNull(product.brand, product.manufacturer).distinct().takeIf { it.isNotEmpty() }?.let {
            Text(it.joinToString(" · "), style = MaterialTheme.typography.titleMedium)
        }
        priceText(product.price, product.currency)?.let { Text(it, style = MaterialTheme.typography.titleMedium) }
        product.description?.let { Text(it, modifier = Modifier.padding(top = 8.dp)) }
        Text("Barcode: ${product.barcode}", style = MaterialTheme.typography.bodySmall, modifier = Modifier.padding(top = 8.dp))

        if (product.isPlaceholder) {
            Card(modifier = Modifier.fillMaxWidth().padding(top = 12.dp)) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text("This product isn't in EveryReview yet", style = MaterialTheme.typography.titleMedium)
                    Text("Know what it is? Add its name, a photo and other details so others can find it.")
                    OutlinedButton(onClick = onAddDetails, modifier = Modifier.fillMaxWidth().padding(top = 8.dp)) {
                        Text("Add product details")
                    }
                }
            }
        }

        Text(ratingSummaryText(product.reviewCount, product.averageRating), modifier = Modifier.padding(top = 12.dp))
        Button(onClick = onWriteReview, modifier = Modifier.fillMaxWidth().padding(top = 8.dp)) {
            Text("Write a review")
        }
        HorizontalDivider(modifier = Modifier.padding(top = 16.dp))
    }
}

private fun ratingSummaryText(reviewCount: Int, averageRating: Double?): String {
    if (reviewCount == 0) return "No reviews yet"
    val reviewsWord = if (reviewCount == 1) "review" else "reviews"
    return "★ %.1f (%d %s)".format(averageRating ?: 0.0, reviewCount, reviewsWord)
}

/** "€2.49" when the currency is a real ISO code, "2.49 XYZ" otherwise, nothing without a price. */
private fun priceText(price: Double?, currency: String?): String? {
    if (price == null) return null
    val currencyInstance = currency?.let { runCatching { Currency.getInstance(it) }.getOrNull() }
    return if (currencyInstance != null) {
        NumberFormat.getCurrencyInstance().apply { this.currency = currencyInstance }.format(price)
    } else {
        listOfNotNull("%.2f".format(price), currency).joinToString(" ")
    }
}

@Composable
private fun ReviewRow(review: Review) {
    Column(modifier = Modifier.fillMaxWidth().padding(vertical = 4.dp)) {
        val ratingText = review.rating?.let { "★".repeat(it) }.orEmpty()
        Text(text = "${review.authorName}  $ratingText", style = MaterialTheme.typography.titleSmall)
        Text(review.body)
    }
}
