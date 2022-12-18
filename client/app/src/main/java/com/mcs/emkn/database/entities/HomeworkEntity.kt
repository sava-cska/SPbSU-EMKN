package com.mcs.emkn.database.entities

import androidx.room.Entity
import androidx.room.PrimaryKey
import com.mcs.emkn.database.entities.HomeworkEntity.Companion.HOMEWORKS_TABLE_NAME
import com.mcs.emkn.network.dto.response.ScoreStatus
import com.mcs.emkn.network.dto.response.TextStatus
import com.squareup.moshi.Json


@Entity(tableName = HOMEWORKS_TABLE_NAME)
data class HomeworkEntity(
    val courseId: Int,
    @PrimaryKey
    val id: Int,
    val name: String,
    val deadline: Long,
    val statusNotPassed: String?,
    val statusUnchecked: String?,
    val score: Int,
    val totalScore: Int
) {
    companion object {
        const val HOMEWORKS_TABLE_NAME = "homeworks"
    }
}