package com.mcs.emkn.network.dto.request

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class UploadImageRequestDto(
    @Json(name = "encoded_jpg") val imageBase64: String,
)