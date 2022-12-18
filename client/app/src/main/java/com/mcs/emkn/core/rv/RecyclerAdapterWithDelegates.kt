package com.mcs.emkn.core.rv

import android.view.ViewGroup
import androidx.recyclerview.widget.RecyclerView

class RecyclerAdapterWithDelegates(
    private val delegates: List<RecyclerDelegate<*, *>>,
    var items: List<Item>,
) : RecyclerView.Adapter<RecyclerView.ViewHolder>() {

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): RecyclerView.ViewHolder {
        return delegates[viewType].onCreateViewHolder(parent)
    }

    override fun onBindViewHolder(holder: RecyclerView.ViewHolder, position: Int) {
        val delegate =
            delegates[getItemViewType(position)] as RecyclerDelegate<ViewHolder<Item>, Item>
        delegate.onBindViewHolder(holder as ViewHolder<Item>, items[position])
    }

    override fun getItemCount(): Int {
        return items.size
    }

    override fun getItemId(position: Int): Long {
        return if (position in items.indices) {
            items[position].getItemId()
        } else
            RecyclerView.NO_ID
    }

    override fun getItemViewType(position: Int): Int {
        val index = delegates.indexOfFirst { delegate -> delegate.matchesItem(items[position]) }
        assert(index >= 0) { "No delegate for item" }
        return index
    }
}