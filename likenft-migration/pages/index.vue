<template>
  <div :class="['flex-1', 'min-h-0', 'bg-likecoin-lightergrey']">
    <div :class="['-mb-[65px]']">
      <HeroBanner>
        <h1
          :class="[
            'text-3xl',
            'font-inter',
            'font-semibold',
            'text-likecoin-votecolor-yes',
          ]"
        >
          {{ $t('app.title') }}
        </h1>
      </HeroBanner>
    </div>
    <div :class="['relative', 'z-10', 'max-w-[880px]', 'px-4', 'mx-auto']">
      <div :class="['bg-white', 'p-[30px]', 'rounded-md', 'shadow-md']">
        <StepSection :step="1" :current-step="currentStep">
          <h2 :class="['text-base', 'font-semibold', 'text-likecoin-darkgrey']">
            {{ $t('section.introduction.title') }}
          </h2>
          <div :class="['mt-2', 'mb-5']">
            <SectionIntroduction
              @confirmClicked="handleIntroductionSectionConfirmClick"
            />
          </div>
        </StepSection>
        <StepSection :step="2" :current-step="currentStep"></StepSection>
        <StepSection :step="3" :current-step="currentStep"></StepSection>
        <StepSection :step="4" :current-step="currentStep"></StepSection>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';

export default Vue.extend({
  computed: {
    currentStep() {
      return 1;
    },
  },
  methods: {
    handleIntroductionSectionConfirmClick() {
      alert('TODO');
    },
  },
});
</script>
