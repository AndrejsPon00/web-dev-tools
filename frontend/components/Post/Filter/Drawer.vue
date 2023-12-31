<template>
    <q-drawer
        v-model="drawerState"
        elevated
        show-if-above
        no-swipe-open
        :width="360"
        side="right"
        class="app-drawer"
    >
        <div class="app-drawer__prepend">
            <q-item>
                <q-item-section avatar>
                    <q-icon :name="mdiFilter" />
                </q-item-section>
                <q-item-section>
                    <q-item-label>Фильтр</q-item-label>
                </q-item-section>
            </q-item>
            <q-separator />
            <div class="q-pa-md">
                <div class="text-body2 q-mb-xs">
                    Поиск:
                </div>
                <post-filter-search-input
                    v-model="query"
                    :readonly="props.readonly"
                    @apply-filters="applyFilters"
                />
            </div>
            <q-separator />
        </div>
        <div class="app-drawer__content">
            <div class="q-pa-md">
                <div class="column q-col-gutter-md">
                    <div class="col">
                        <post-filter-source-select
                            v-model="sources"
                            :sources="filter.sources"
                            :loading="filterDataLoading"
                            :readonly="props.readonly"
                        />
                    </div>
                    <div class="col">
                        <post-filter-price-range
                            v-model="price"
                            :readonly="props.readonly"
                        />
                    </div>
                </div>
            </div> 
        </div>
        <div class="app-drawer__append">
            <q-separator />
            <div class="q-pa-sm">
                <div class="row q-col-gutter-sm">
                    <div class="col-auto">
                        <q-btn
                            :icon="mdiUpdate"
                            :disable="props.readonly"
                            color="accent"
                            @click.passive="clearFilters"
                        >
                            <q-tooltip
                                anchor="top middle"
                                self="bottom middle"
                            >
                                Сбросить
                            </q-tooltip>
                        </q-btn>
                    </div>
                    <div class="col-grow">
                        <q-btn
                            class="full-width"
                            color="primary"
                            :loading="props.readonly"
                            @click.passive="applyFilters"
                        >
                            <q-icon
                                left
                                :name="mdiLayersSearch"
                            />
                            Применить
                        </q-btn>
                    </div>
                </div>
            </div>
        </div>
    </q-drawer>
</template>

<script lang="ts" setup>
import { mdiFilter, mdiLayersSearch, mdiUpdate } from '@quasar/extras/mdi-v7';
import { PriceRange } from './PriceRange.vue';

export type FilterFields = {
    query: string | undefined;
    categories: string[] | string | undefined;
    sources: string[] | string | undefined;
    price: {
        min: number | undefined;
        max: number | undefined;
    },
};

export type Props = {
    modelValue: FilterFields;
    readonly?: boolean;
};

export type Emits = {
    (e: 'update:modelValue', modelValue: FilterFields): void;
};

const emits = defineEmits<Emits>();
const props = withDefaults(defineProps<Props>(), {
    readonly: false,
});

const query = ref<string | undefined>(props.modelValue.query);
const sources = ref<string[] | string | undefined>(props.modelValue.sources);
const categories = ref<string[] | string | undefined>(props.modelValue.categories);
const price = ref<PriceRange>({ min: props.modelValue.price.min, max: props.modelValue.price.max });

const filterData = useFilterData();
const drawerState = useFilterDrawerState();

const applyFilters = () => {
    emits('update:modelValue', {
        query: query.value,
        categories: categories.value,
        sources: sources.value,
        price: {
            max: price.value.max,
            min: price.value.min,
        },
    });
};

const clearFilters = () => {
    query.value = undefined;
    sources.value = [];
    categories.value = [];
    price.value = { min: undefined, max: undefined };

    applyFilters();
};

watch(() => props.modelValue, (newModelValue) => {
    query.value = newModelValue.query;
    sources.value = newModelValue.sources;
    categories.value = newModelValue.categories;
    price.value = {
        min: newModelValue.price.min,
        max: newModelValue.price.max,
    };
});

const filterSources = useFilterSourcesState();

const { data: filter, pending: filterDataLoading } = useAsyncData('filter-data', () => filterData.fetchFilterData(), {
    default: () => ({
        sources: [],
    }),
    lazy: true,
});

watch(filter, (newFilter) => {
    filterSources.value = newFilter.sources.map((v) => v.value);
});
</script>