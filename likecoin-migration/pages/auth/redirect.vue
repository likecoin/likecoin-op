<script>
export default {
  fetch({ redirect, query }) {
    const { code, method, ...others } = query;
    const queries = new URLSearchParams({
      code,
      method,
      ...others,
    });
    redirect(`/extended?${queries.toString()}`);
  },
};
</script>
