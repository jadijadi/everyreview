package app.everyreview.data.local

import androidx.room.Database
import androidx.room.RoomDatabase
import app.everyreview.data.local.entity.ProductEntity
import app.everyreview.data.local.entity.ReviewEntity

@Database(entities = [ProductEntity::class, ReviewEntity::class], version = 2, exportSchema = false)
abstract class EveryReviewDatabase : RoomDatabase() {
    abstract fun productDao(): ProductDao

    abstract fun reviewDao(): ReviewDao
}
