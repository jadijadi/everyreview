package app.everyreview.di

import android.content.Context
import androidx.room.Room
import app.everyreview.data.local.EveryReviewDatabase
import app.everyreview.data.local.ProductDao
import app.everyreview.data.local.ReviewDao
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object DatabaseModule {

    @Provides
    @Singleton
    fun provideDatabase(@ApplicationContext context: Context): EveryReviewDatabase =
        Room.databaseBuilder(context, EveryReviewDatabase::class.java, "everyreview.db").build()

    @Provides
    fun provideProductDao(database: EveryReviewDatabase): ProductDao = database.productDao()

    @Provides
    fun provideReviewDao(database: EveryReviewDatabase): ReviewDao = database.reviewDao()
}
