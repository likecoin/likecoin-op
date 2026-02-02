<template>
  <div :class="['flex-1', 'min-h-0', 'bg-likecoin-lightergrey', 'pb-4']">
    <div :class="['-mb-[65px]']">
      <HeroBanner>
        <h1
          :class="[
            'px-4',
            'text-2xl',
            'sm:text-3xl',
            'font-inter',
            'font-semibold',
            'text-likecoin-votecolor-yes',
          ]"
        >
          {{ $t('app.title') }}
        </h1>
      </HeroBanner>
    </div>
    <div :class="['relative', 'max-w-[880px]', 'sm:px-4', 'mx-auto']">
      <div
        :class="[
          'bg-white',
          'p-[30px]',
          'max-sm:px-4',
          'sm:rounded-md',
          'shadow-md',
          'flex',
          'flex-col',
          'items-center',
          'gap-4',
        ]"
      >
        <img src="~/assets/images/icon.png" :class="['w-20', 'h-auto']" />
        <div :class="['mt-2', 'mb-5']">
          <SectionEnded />
        </div>
      </div>
    </div>
  </div>
</template>
