package com.mcs.emkn.core

import android.content.Context
import android.view.View
import android.view.ViewGroup

interface Component {
    fun createView(context: Context, parent: ViewGroup?): View
    
    fun onViewCreated(view: View) = Unit
}