package app.everyreview.domain.usecase

import app.everyreview.domain.model.Product
import app.everyreview.domain.model.ProductDetails
import app.everyreview.domain.repository.ProductRepository
import javax.inject.Inject

class SubmitProductDetailsUseCase @Inject constructor(
    private val repository: ProductRepository,
) {
    suspend operator fun invoke(productId: String, details: ProductDetails): Result<Product> {
        val name = details.name.trim()
        if (name.isEmpty()) {
            return Result.failure(IllegalArgumentException("Product name must not be empty"))
        }
        val price = details.price?.trim()?.takeIf { it.isNotEmpty() }
        if (price != null && price.replace(',', '.').toDoubleOrNull()?.takeIf { it >= 0 } == null) {
            return Result.failure(IllegalArgumentException("Price must be a non-negative number"))
        }
        return repository.submitDetails(
            productId,
            ProductDetails(
                name = name,
                brand = details.brand.orBlankToNull(),
                manufacturer = details.manufacturer.orBlankToNull(),
                description = details.description.orBlankToNull(),
                price = price,
                currency = details.currency.orBlankToNull()?.uppercase(),
                photo = details.photo,
            ),
        )
    }

    private fun String?.orBlankToNull(): String? = this?.trim()?.takeIf { it.isNotEmpty() }
}
