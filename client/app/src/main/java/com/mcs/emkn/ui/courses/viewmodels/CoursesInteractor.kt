package com.mcs.emkn.ui.courses.viewmodels

import com.mcs.emkn.core.State
import com.mcs.emkn.database.entities.CheckBoxesStateEntity
import com.mcs.emkn.database.entities.PeriodEntity
import kotlinx.coroutines.Deferred
import kotlinx.coroutines.flow.Flow

interface CoursesInteractor {
    val courses: Flow<State<PeriodCourses>>

    val periods: Flow<PeriodsData>

    val period: Flow<PeriodEntity>

    val navEventsFlow: Flow<CoursesNavEvents>

    val errorsFlow: Flow<CoursesError>

    fun onPeriodChosen(periodId: Int)

    fun loadPeriodsAndCourses()

    fun changeEnrollState(id: Int)
}