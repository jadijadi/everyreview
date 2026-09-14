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
    val manufacturer: String? = null,
    val description: String? = null,
    val price: Double? = null,
    val currency: String? = null,
)
