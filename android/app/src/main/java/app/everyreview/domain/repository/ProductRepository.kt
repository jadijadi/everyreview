package app.everyreview.domain.repository

import app.everyreview.domain.model.Product

interface ProductRepository {
    suspend fun getByBarcode(barcode: String): Result<Product>
}
