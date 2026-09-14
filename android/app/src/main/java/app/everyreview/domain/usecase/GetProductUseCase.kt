package app.everyreview.domain.usecase

import app.everyreview.domain.model.Product
import app.everyreview.domain.repository.ProductRepository
import javax.inject.Inject

class GetProductUseCase @Inject constructor(
    private val repository: ProductRepository,
) {
    suspend operator fun invoke(barcode: String): Result<Product> = repository.getByBarcode(barcode)
}
