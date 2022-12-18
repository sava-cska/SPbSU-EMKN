package com.mcs.emkn.database

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import com.mcs.emkn.database.entities.*


@Dao
interface CoursesDao{
    @Query("SELECT * FROM ${CourseEntity.COURSES_TABLE_NAME}")
    fun getAllCourses(): List<CourseEntity>

    @Query("SELECT * FROM ${CourseEntity.COURSES_TABLE_NAME} WHERE periodId IN (:ids)")
    fun getCoursesByPeriods(ids: List<Int>): List<CourseEntity>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    fun putCourses(courses: List<CourseEntity>)

    @Query("UPDATE ${CourseEntity.COURSES_TABLE_NAME} SET enrolled=:enrolled WHERE id = :id")
    fun updateCourseEnrolledState(id: Int, enrolled: Boolean)

    @Query("DELETE FROM ${CourseEntity.COURSES_TABLE_NAME}")
    fun deleteCourses()

    @Query("SELECT * FROM ${PeriodEntity.PERIODS_TABLE_NAME}")
    fun getPeriods(): List<PeriodEntity>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    fun putPeriods(periods: List<PeriodEntity>)

    @Query("DELETE FROM ${PeriodEntity.PERIODS_TABLE_NAME}")
    fun deletePeriods()

    @Query("SELECT * FROM ${ProfileEntity.PROFILES_TABLE_NAME}")
    fun getProfiles(): List<ProfileEntity>

    @Query("SELECT * FROM ${ProfileEntity.PROFILES_TABLE_NAME} WHERE id IN (:ids)")
    fun getProfilesByIds(ids: List<Int>): List<ProfileEntity>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    fun putProfiles(profiles: List<ProfileEntity>)

    @Query("DELETE FROM ${ProfileEntity.PROFILES_TABLE_NAME}")
    fun deleteProfiles()

    @Query("SELECT * FROM ${CheckBoxesStateEntity.CHECK_BOX_TABLE_NAME}")
    fun getCheckBoxState(): List<CheckBoxesStateEntity>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    fun putCheckBoxState(checkBoxStateEntity: CheckBoxesStateEntity)

    @Query("DELETE FROM ${CheckBoxesStateEntity.CHECK_BOX_TABLE_NAME}")
    fun deleteCheckBoxState()

    @Query("SELECT * FROM ${HomeworkEntity.HOMEWORKS_TABLE_NAME} WHERE courseId = :id")
    fun getHomeworksByCourse(id: Int): List<HomeworkEntity>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    fun putHomeworks(homeworks: List<HomeworkEntity>)
}