package app.everyreview.data.local.entity

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "products")
data class ProductEntity(
    @PrimaryKey val barcode: String,
    val id: String,
    val name: String?,
    val brand: String?,
    val imageUrl: String?,
    val isPlaceholder: Boolean,
    val reviewCount: Int,
    val averageRating: Double?,
    val manufacturer: String?,
    val description: String?,
    val price: Double?,
    val currency: String?,
    val cachedAt: Long,
)
