package app.everyreview.domain.repository

import app.everyreview.domain.model.Product
import app.everyreview.domain.model.ProductDetails

interface ProductRepository {
    suspend fun getByBarcode(barcode: String): Result<Product>

    suspend fun submitDetails(productId: String, details: ProductDetails): Result<Product>
}
