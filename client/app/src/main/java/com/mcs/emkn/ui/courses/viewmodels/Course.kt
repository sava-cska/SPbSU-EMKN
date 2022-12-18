package com.mcs.emkn.ui.courses.viewmodels

import com.mcs.emkn.database.entities.ProfileEntity
import com.squareup.moshi.Json

data class Course(
    val id: Int,
    val title: String,
    val enrolled: Boolean,
    val shortDescription: String,
    val teachersProfiles: List<Int>
)
