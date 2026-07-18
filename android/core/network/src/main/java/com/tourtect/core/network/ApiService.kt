package com.tourtect.core.network

import com.tourtect.core.model.Place
import com.tourtect.core.model.Post
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.Path
import retrofit2.http.Query

interface ApiService {
    @GET("v1/places")
    suspend fun getPlaces(
        @Query("q") query: String? = null,
        @Query("category") category: String? = null
    ): ListResponse<Place>

    @GET("v1/places/{placeId}")
    suspend fun getPlaceDetail(
        @Path("placeId") placeId: String
    ): Place

    @GET("v1/posts")
    suspend fun getPublishedPosts(
        @Query("place_id") placeId: String? = null,
        @Query("post_type") postType: String? = null
    ): ListResponse<Post>

    @POST("v1/posts/drafts")
    suspend fun createDraft(
        @Body request: CreateDraftRequest
    ): Post

    @POST("v1/posts/{postId}/publish")
    suspend fun publishPost(
        @Path("postId") postId: String
    ): Post
}

data class ListResponse<T>(
    val items: List<T>,
    val pagination: PaginationInfo
)

data class PaginationInfo(
    val nextCursor: String?,
    val hasMore: Boolean
)

data class CreateDraftRequest(
    val post_type: String,
    val original_locale: String,
    val title: String,
    val body: String,
    val place_ids: List<String>? = null
)
