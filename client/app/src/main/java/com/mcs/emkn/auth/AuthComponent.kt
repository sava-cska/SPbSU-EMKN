package com.mcs.emkn.auth

import android.content.Context
import android.view.View
import android.view.ViewGroup
import com.mcs.emkn.core.Component
import com.mcs.emkn.core.Router
import com.mcs.emkn.database.Database
import com.mcs.emkn.network.Api
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class AuthComponent @Inject constructor(): Component {
    @Inject
    lateinit var router: Router
    @Inject
    lateinit var api: Api
    @Inject
    lateinit var database: Database
    
    override fun createView(context: Context, parent: ViewGroup?): View {
        return View(context)
    }
}