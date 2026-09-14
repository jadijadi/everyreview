package app.everyreview.domain.model

data class Product(
    val id: String,
    val barcode: String,
    val name: String?,
    val brand: String?,
    val imageUrl: String?,
    val isPlaceholder: Boolean,
    val reviewCount: Int,
    val averageRating: Double?,
)
