package app.everyreview.domain.model

import java.io.File

/**
 * What a user fills in for a product nobody has described yet. Only [name] is required;
 * blank optional fields are sent as absent. [photo] is an already-downscaled JPEG on disk.
 */
data class ProductDetails(
    val name: String,
    val brand: String? = null,
    val manufacturer: String? = null,
    val description: String? = null,
    val price: String? = null,
    val currency: String? = null,
    val photo: File? = null,
)
