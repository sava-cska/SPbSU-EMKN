package com.mcs.emkn.core.rv

import android.view.ViewGroup

interface RecyclerDelegate<VH : ViewHolder<I>, I : Item> {
    fun onBindViewHolder(viewHolder: VH, item: I)

    fun onCreateViewHolder(parent: ViewGroup): VH

    fun matchesItem(item: Item): Boolean
}