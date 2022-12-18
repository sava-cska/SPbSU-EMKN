package com.mcs.emkn.core.rv

import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.annotation.LayoutRes
import androidx.recyclerview.widget.RecyclerView

abstract class ViewHolder<I : Item>(view: View) : RecyclerView.ViewHolder(view) {
    abstract fun bind(item: I)

    constructor(@LayoutRes layoutRes: Int, parent: ViewGroup) : this(
        LayoutInflater.from(parent.context).inflate(layoutRes, parent, false)
    )
}