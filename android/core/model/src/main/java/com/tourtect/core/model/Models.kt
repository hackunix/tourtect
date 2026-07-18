package com.tourtect.core.model

import java.util.Date

data class Place(
    val placeId: String,
    val name: String,
    val category: String,
    val regionId: String,
    val address: String?,
    val latitude: Double,
    val longitude: Double,
    val postCount: Int,
    val averageRating: Double,
    val freshness: Date?
)

data class Post(
    val postId: String,
    val authorId: String,
    val postType: String,
    val originalLocale: String,
    val title: String,
    val body: String,
    val evidenceLevel: String,
    val commercialDisclosure: String,
    val moderationStatus: String,
    val createdAt: Date
)
