package app.everyreview.data.local.entity

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "reviews")
data class ReviewEntity(
    @PrimaryKey val id: String,
    val productId: String,
    val authorName: String,
    val body: String,
    val rating: Int?,
    val createdAt: String,
)
