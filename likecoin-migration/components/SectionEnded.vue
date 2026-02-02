<template>
  <div :class="['p-4', 'rounded-xl', 'border', 'border-likecoin-grey']">
    <i18n
      path="section.ended.text"
      tag="p"
      :class="['text-sm', 'text-likecoin-darkgrey', 'whitespace-pre-line']"
    >
      <template #link>
        <a
          href="https://likecoin.substack.com/8th-anniversary"
          target="_blank"
          :class="['text-like-green', 'underline']"
          >{{ $t('section.ended.substack-link-text') }}</a
        >
      </template>
    </i18n>
  </div>
</template>
