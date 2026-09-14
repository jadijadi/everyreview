package app.everyreview.data.repository

import app.everyreview.data.local.entity.ProductEntity
import app.everyreview.data.local.entity.ReviewEntity
import app.everyreview.data.remote.dto.ProductResponse
import app.everyreview.data.remote.dto.ReviewResponse
import app.everyreview.domain.model.Product
import app.everyreview.domain.model.Review

fun ProductResponse.toDomain(): Product = Product(
    id = id,
    barcode = barcode,
    name = name,
    brand = brand,
    imageUrl = imageUrl,
    isPlaceholder = isPlaceholder,
    reviewCount = ratingSummary.reviewCount,
    averageRating = ratingSummary.averageRating,
)

fun Product.toEntity(cachedAt: Long = System.currentTimeMillis()): ProductEntity = ProductEntity(
    barcode = barcode,
    id = id,
    name = name,
    brand = brand,
    imageUrl = imageUrl,
    isPlaceholder = isPlaceholder,
    reviewCount = reviewCount,
    averageRating = averageRating,
    cachedAt = cachedAt,
)

fun ProductEntity.toDomain(): Product = Product(
    id = id,
    barcode = barcode,
    name = name,
    brand = brand,
    imageUrl = imageUrl,
    isPlaceholder = isPlaceholder,
    reviewCount = reviewCount,
    averageRating = averageRating,
)

fun ReviewResponse.toDomain(): Review = Review(
    id = id,
    productId = productId,
    authorName = authorName,
    body = body,
    rating = rating,
    createdAt = createdAt,
)

fun Review.toEntity(): ReviewEntity = ReviewEntity(
    id = id,
    productId = productId,
    authorName = authorName,
    body = body,
    rating = rating,
    createdAt = createdAt,
)

fun ReviewEntity.toDomain(): Review = Review(
    id = id,
    productId = productId,
    authorName = authorName,
    body = body,
    rating = rating,
    createdAt = createdAt,
)
