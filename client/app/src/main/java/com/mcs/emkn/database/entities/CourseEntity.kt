package com.mcs.emkn.database.entities

import androidx.room.Entity
import androidx.room.PrimaryKey
import com.mcs.emkn.database.entities.CourseEntity.Companion.COURSES_TABLE_NAME

@Entity(tableName = COURSES_TABLE_NAME)
data class CourseEntity(
    val periodId: Int,
    @PrimaryKey
    val id: Int,
    val title: String,
    var enrolled: Boolean,
    val shortDescription: String,
    val teachersProfiles: List<Int>
) {
    companion object {
        const val COURSES_TABLE_NAME = "courses"
    }
}