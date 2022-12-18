package com.mcs.emkn.database.entities

import androidx.room.Entity
import androidx.room.PrimaryKey
import com.mcs.emkn.database.entities.ProfileEntity.Companion.PROFILES_TABLE_NAME
import com.squareup.moshi.Json


@Entity(tableName = PROFILES_TABLE_NAME)
data class ProfileEntity(
    @PrimaryKey
    val id: Int,
    val avatarUrl: String,
    val firstName: String,
    val secondName: String
) {
    companion object {
        const val PROFILES_TABLE_NAME = "profiles"
    }
}


