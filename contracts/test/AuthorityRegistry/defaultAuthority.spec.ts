describe('AuthorityRegistry Test', function () {
  describe('DefaultAuthority', function () {
    describe('unhappy paths', function () {
      it('reverts when non-owner adds an authority', async () => {});
      it('reverts when non-owner removes an authority', async () => {});
      it('reverts when adding an already-present authority', async () => {});
      it('reverts when removing a non-present authority', async () => {});
    });

    describe('happy paths', function () {
      it('allows owner to add authority and registry reflects it', async () => {});
      it('allows owner to remove authority and registry reflects it', async () => {});
    });
  });
});
